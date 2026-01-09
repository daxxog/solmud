We are currently working on getting the 317 deob located in ./srcAllDummysRemoved mapped 1:1 with the bytecode in rs317og.jar (extracted to ./bytecode/)


73 classes, excluded GUI
```
% ls -lah srcAllDummysRemoved/.javap_cache | grep -v GUI | grep '.javap.cache' | wc -l
      73
```


73 classes (OG)
```
% tree bytecode/client | grep files
1 directory, 73 files
```

Probably correct mapping (unless proven WRONG using **evidence**):
`bytecode/mapping/class_mapping.csv`


# The goal is to use an evidence-based approach to prove the mapping in `bytecode/mapping`.


## Each evidence file should:
```
_note for agent_
You can break the research needed to complete these tasks as delegations for subagents.
```

 - [ ] Contain an overview of the class, it's intended purpose, and core functionality.
 - [ ] Explain its role in the overall architecture and relationship with other classes. Use mermaid diagrams where appropriate for documenting class-relationships.
 - [ ] Contain bash commands which can be ran to show specific areas of bytecode that prove a mapping.
 - [ ] Have commands that show sections of DEOB source code that relate to each section of obfusticated bytecode used as evidence.
 - [ ] Have commands that show sections of DEOB javap cache that relate to each section of obfusticated bytecode used as evidence.
 - [ ] Not just grep for a single line—multiple lines of evidence should be shown to display the context around the bytecode segment.
 - [ ] Be verified by checking the commands actually work and the presented evidence is non-contradictory.
 - [ ] Only be mapped 1:1—each obfusticated bytecode file should only map to a single DEOB source file/javap cache.
 - [ ] **NOT** use absolute paths in documentation (such as `/Users/daxxog/Desktop`). Paths should be relative to the project root.



For any disputes/contradictions document the issue in `bytecode/mapping/evidence/disputed`, then use subagents to research what the proper mapping _should_ be.
