We are currently working on getting the 317 deob located in ./srcAllDummysRemoved mapped 1:1 with the bytecode in rs317og.jar (extracted to ./bytecode/). Don't trust—verify. If something doesn't seem right (event the csv!), create a comprehensive plan to fix it. This document is the single source of truth for our task. This specific document **must not** be modified by any agent for any reason.


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

    - Make sure each subagent **reads** this file (./OG_vs_DEOB.md) _before_ editing any evidence files!
    - Run jobs in batches to avoid rate limits.
    - Tell the agent they are a subagent.

    _ If you are _not_ a subagent reading this, your *primary task* is to verify _each and every_ mapping file to adhere to this document, using the subagent strategy documented above! Don't trust—verify!
```

### ADHERE TO CRITICAL CHECKLIST
 - [ ] If you have never analyzed this before (and are not a subagent)—Rank evidence by quality and use the highest quality set of evidence as the template. Existing "template" files should be ignored and removed. We only want exceptional forensic-grade evidence with beautiful documentation.
 - [ ] Diagrams should only contain reference to classes that exist in the DEOB source! Only reference classes by their DEOB name in diagrams (diagrams full of OG names are _useless_).
 - [ ] Contain an overview of the class, it's intended purpose, and core functionality.
 - [ ] Explain its role in the overall architecture and relationship with other classes. Use mermaid diagrams where appropriate for documenting class-relationships (inheritance, data access, method calls).
 - [ ] Contain bash commands which can be ran to show specific areas of bytecode that prove a mapping.
 - [ ] Have commands that show sections of DEOB source code that relate to each section of obfusticated bytecode used as evidence.
 - [ ] Have commands that show sections of DEOB javap cache that relate to each section of obfusticated bytecode used as evidence.
 - [ ] Not just grep for a single line—multiple lines of evidence should be shown to display the context around the bytecode segment.
 - [ ] Be verified by checking the commands actually work and the presented evidence is non-contradictory.
 - [ ] Only be mapped 1:1—each obfusticated bytecode file should only map to a single DEOB source file/javap cache.
 - [ ] **NOT** use absolute paths in documentation (such as `/Users/daxxog/Desktop`). Paths should be relative to the project root.
 - [ ] `bytecode/mapping` should be a clean, auditable directory without unnecessary files, should not contain READMES unaligned to this document, random "docs" or "templates", and markdown spam with notes from previous agentic runs.
 - [ ] If disputes are resolved they are no longer disputes. So there shouldn't be "resolved disputes" documented anywhere.



For any disputes/contradictions document the issue in `bytecode/mapping/evidence/disputed`, then use subagents to research what the proper mapping _should_ be.


### Notes to avoid wrong-direction and rabbit trails
 - [ ]
 - Quantity != Quality! Counting the number of lines in a file does not imply those lines are of quality craftsmanship.
 - You can ignore the old tool in `./tools/classmapper_ARCHIVED_2026-01-08/`. This is archived and no longer apart of the strategy.
 - Avoid creating files like `bytecode/mapping/evidence/verified/Animable.md`. This doesn't contain any context in the filename about _which_ OG class it is mapped to! Instead use a path like `bytecode/mapping/evidence/verified/RSSocket_NQABEVLK.md`. Notice how the case is preserved.
 - Correct order for file naming convention is DEOB (first) _ OG (second)
 - OG = ./rs317og.jar (you can list the contents of the jar if you are confused)
 - DEOB = ./srcAllDummysRemoved.zip (you can list the contents of the zip if you are confused)
 - create a `./Makefile` with helper commands that speed up this work when it makes sense. Modify existing `./Makefile` if it does not align to this document.
 - "Build mode" does NOT mean creating a comprehensive build system. It simply refers to "read-write-mode" of the agent (i.e. !PLAN_MODE);
 - Random `_archive` or `backup` directories should be removed. (see requirement on keeping `bytecode/mapping` clean)
 - You are only responsible for managing `./Makefile` and the `bytecode/mapping` folder. Don't concern yourself with the rest of this project at this stage. Only read paths relevant to the goals outlined in this file.
 - Focus on solving CSV-completeness, accuracy, and dispute resolution first. Then laser-focus on evidence quality.
 - Don't ask the user about priority. The mantra describes priority already and contains a clear definition of what good looks like.
 - `signlink` does not need a mapping, however its relationship to other classes may be of interest in relevant mermaid diagrams.
 - If you are thinking about saying something like "Ready to execute the verification protocol.", _just do it_™️. Multi-pass and use subagents to reach deep-planning goals.
