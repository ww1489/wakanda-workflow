#!/usr/bin/env python3
"""Count versioned text including new source files, excluding ignored artifacts."""
import argparse,collections,json,pathlib,re,subprocess,tarfile
p=argparse.ArgumentParser();p.add_argument('--ref');p.add_argument('--archive');p.add_argument('--output');a=p.parse_args()
root=pathlib.Path(subprocess.check_output(['git','rev-parse','--show-toplevel'],text=True).strip())
exts={'.go','.ts','.tsx','.js','.jsx','.mjs','.cjs','.py','.sh','.html','.css','.less','.scss','.sass','.sql','.vue','.svelte'}
counts=collections.defaultdict(lambda:dict(files=0,lines=0,bytes=0));total=dict(files=0,lines=0,bytes=0); manifests=0;rush=0;godeps=0
proc=None
archive=None
all_files=0
front_declarations=0
front_external=set()
front_workspace_edges=0
lock_entries=0
if a.ref:
 items=[]
 for x in subprocess.check_output(['git','ls-tree','-rz',a.ref]).split(b'\0'):
  if x:
   meta,path=x.split(b'\t',1);mode,kind,sha=meta.decode().split()
   if kind=='blob':items.append((path.decode(),sha))
 proc=subprocess.Popen(['git','cat-file','--batch'],stdin=subprocess.PIPE,stdout=subprocess.PIPE)
elif a.archive:
 archive=tarfile.open(a.archive);items=[(m.name.removeprefix('./'),m) for m in archive.getmembers() if m.isfile() or m.issym()]
else:items=[(x.decode(),None) for x in sorted(set(subprocess.check_output(['git','ls-files','-co','--exclude-standard','-z']).split(b'\0'))- {b''})]
for path,sha in items:
 if archive:
  data=sha.linkname.encode() if sha.issym() else archive.extractfile(sha).read()
 elif proc:
  proc.stdin.write(sha.encode()+b'\n');proc.stdin.flush();n=int(proc.stdout.readline().split()[2]);data=proc.stdout.read(n);proc.stdout.read(1)
 else:
  f=root/path
  if not f.exists():continue
  data=f.readlink().as_posix().encode() if f.is_symlink() else f.read_bytes()
 all_files+=1
 if b'\0' in data:continue
 try:text=data.decode('utf-8')
 except UnicodeDecodeError:continue
 f=pathlib.Path(path);test=f.name.endswith('_test.go') or '.test.' in f.name or '.spec.' in f.name or '__tests__' in f.parts or 'tests' in f.parts
 generated=bool(re.search(r'(?im)^\s*(?://|/\*|\*)?\s*(?:Code generated|This file is (?:auto|automatically)[ -]generated|Automatically generated)',text[:2000])) or 'auto-generated' in f.parts or f.name.endswith('.gen.go')
 if f.suffix in exts:
  prefix='backend' if path.startswith('backend/') else 'frontend' if path.startswith('frontend/') else 'other'
  group=prefix+('_tests' if test else '_generated' if generated else '_source')
 elif f.suffix=='.md':group='docs'
 elif path.startswith(('docker/','helm/','scripts/')) or f.name in ('Makefile','Dockerfile'):group='deployment_config'
 else:group='other_text'
 values=dict(files=1,lines=data.count(b'\n')+int(bool(data) and not data.endswith(b'\n')),bytes=len(data))
 for k,v in values.items():counts[group][k]+=v;total[k]+=v
 if f.name=='package.json':
  manifests+=1
  if path.startswith('frontend/'):
   manifest=json.loads(text)
   for dep,version in manifest.get('dependencies',{}).items():
    front_declarations+=1
    if str(version).startswith('workspace:'):front_workspace_edges+=1
    else:front_external.add(dep)
 if path.startswith('common/config/subspaces/') and f.name=='pnpm-lock.yaml':lock_entries+=len(re.findall(r'^  /[^\n]+:',text,re.M))
 if path=='rush.json':rush=len(re.findall(r'"packageName"\s*:',text))
 if path=='backend/go.mod':godeps=len(re.findall(r'^\s+\S+\s+v\S+',text,re.M))
if proc:proc.stdin.close();proc.wait()
result=dict(ref=a.ref or a.archive or 'working-tree',total=total,categories=counts,package_manifests=manifests,rush_projects=rush,go_requirements=godeps,all_files=all_files,frontend_production_dependency_declarations=front_declarations,frontend_external_production_dependencies=len(front_external),frontend_workspace_dependency_edges=front_workspace_edges,pnpm_locked_package_entries=lock_entries)
encoded=json.dumps(result,ensure_ascii=False,indent=2)
if a.output:pathlib.Path(a.output).write_text(encoded+'\n')
print(encoded)
