# The Cartographer of Lost Things

## A Short Story

I found it on a Tuesday afternoon, which seems appropriate. Tuesdays are in-between days—not the fresh promise of Monday, not the downhill slide toward Friday. Just existing, like most of the internet.

I was looking for something to fix. Not in the grand, world-saving sense, but in the small, satisfying way of finding a bug, understanding it, and making it right. My therapist says I need projects with clear endpoints. She doesn't know about this yet.

The repository was called `solmud`. A friend from the vintage computing forum mentioned it in passing—"someone's trying to resurrect the old RuneScape 317 client." I had fond memories of that era, those blocky characters and tinny MIDI music from 2005. I was thirteen then, and the world felt simpler, or perhaps I was just simpler in it.

I clicked.

The first thing I noticed was that it wasn't just a game server. That's what I expected—`main.go`, some protocol handlers, maybe a README with build instructions. But there were *categories*. Directories like archaeological strata.

```
solmud/
├── solmud_server/
│   ├── main.go
│   └── cache extraction logic
├── bytecode/
│   └── mapping/
│       ├── evidence/
│       │   ├── verified/
│       │   └── disputed/
│       └── class_mapping.csv
├── eng-web_vpl.zip
└── supercharger/
    └── coords.go
```

Wait. Back up.

`eng-web_vpl.zip`? I opened it. The World English Bible. Complete, in SQL, XML, and plain text. 2020 stable text edition. Five megabytes of sacred scripture, quietly sitting beside obfuscated Java bytecode from twenty years ago.

And `supercharger/coords.go`? Geographic coordinate interfaces. Degrees, minutes, seconds. The mathematics of position on a spherical earth.

I sat back from my monitor. The afternoon light had shifted, becoming that particular shade of gold that makes dust motes visible in sunbeams. I had entered someone's mind, or perhaps their room. A room where a 2005 video game, biblical texts, and cartographic mathematics coexist without apparent contradiction.

I started reading.

---

The `bytecode/mapping/` directory was a forensic investigation. Whoever had created this wasn't just preserving a game; they were *proving* things. There were 73 classes mapped between `rs317og.jar` (original, obfuscated) and `srcAllDummysRemoved` (deobfuscated). Each mapping had evidence. Files in `verified/` with names like `client.md` and `WorldController.md` contained detailed analysis—bytecode patterns cross-referenced with source code structures, Mermaid diagrams showing inheritance chains, confidence scores like 94.0% or 0.95/1.0.

It was obsessive in the way that restoration work is obsessive. The kind of person who removes varnish from a Renaissance painting one square centimeter at a time. Who documents every decision. Who treats the past with the reverence usually reserved for religious artifacts.

And then I remembered: there were religious artifacts here too.

---

I opened `OG_vs_DEOB.md`. It was a master document for the mapping verification system, but the language struck me as almost liturgical in its precision:

> "Every mapping requires forensic-grade evidence. Three-way correlation: bytecode + source + javap cache. Binary PASS/FAIL quality gates with 6 verification criteria."

This wasn't code documentation. This was a *catechism*.

I found `MUSIC_MYSTERY.md` next. It detailed RuneScape 317's unique audio system—how the game used JavaScript LiveConnect to play MIDI files through the browser's external audio capabilities in 2005. The technical explanation was thorough, but the tone was something else. Wonder, perhaps. Or mourning. The author clearly found this obsolete technology beautiful in its ingenuity, in its now-irrelevant cleverness.

I kept thinking about curation. About the choices we make when we decide what to save.

---

The server code in `solmud_server/main.go` was elegant. Modern Go, proper error handling, clean abstractions. Someone had taken a protocol from 2005—on-demand file serving with 520-byte sector-chained cache format—and implemented it in a language that wouldn't exist until 2009. The anachronism felt deliberate, like rebuilding a cathedral using modern engineering but medieval plans.

There were comments about CRC verification issues. Debug logs. The story of someone troubleshooting why a twenty-year-old client wouldn't connect to their fresh server. The frustration and eventual triumph were visible in the commit history, which I pulled up next.

Late December 2024: initial setup, adding original JARs.

January 2025: branching into experiments—`psyop`, `poc1`, `a_bunch_of_random_stuff`.

The branch names told their own story. Someone having fun. Someone treating this seriously but not solemnly. Someone who called their experimentation branch `a_bunch_of_random_stuff` because why pretend to be more organized than you are?

---

I looked at the Bible data again. `eng-web_vpl`—World English Bible, Verse Per Line format. The 2020 stable text edition.

Why?

I searched for references to it in the code. Found nothing direct. It wasn't being used by the server. Wasn't parsed by any tool. Just... present. Five megabytes of text describing creation, fall, redemption, apocalypse. Sitting in a directory with obfuscated Java class files named things like `AFCKELYG` and `ARZPHHDH`.

I thought about curation again. About the medieval monks who copied texts—not just religious ones, but pagan philosophy, scientific treatises, bawdy poems. The monastery as a kind of memory palace, preserving what the world wasn't ready to remember. A lot of what we know about antiquity exists because someone in a scriptorium decided to spend months copying it by hand, not because they agreed with it, but because they recognized that it mattered.

Was this repository a digital scriptorium?

---

I opened `supercharger/coords.go`. It defined interfaces for geographic coordinates—DMS (Degrees Minutes Seconds) notation. The code was clean, generic, type-safe. A small library for converting between coordinate systems, for calculating distances on a sphere.

No obvious connection to RuneScape. The game had a coordinate system, but it was fictional—x and y positions in a virtual world, not latitudes and longitudes. And no connection to the Bible, unless you were plotting the locations of biblical events, which didn't seem to be happening here.

Three domains. Three mapping projects:

1. Mapping obfuscated code to its true names
2. Mapping sacred text to digital formats
3. Mapping the Earth's surface to mathematical abstractions

All different. All the same.

---

There's a concept in cartography called "generalization"—the process of deciding what to show and what to omit when representing the world at different scales. A map at 1:1 would be useless; it would be the territory itself. Every map is a lie that tells a truth through omission.

I think this repository was someone's personal map. Not of geography, but of interest. Of what mattered enough to preserve, document, understand.

The RuneScape 317 client wasn't particularly important in the grand scheme. A mass-market video game from 2005, obsolete technology, abandoned code. But to someone, it was worth months of forensic analysis. Worth building a server from scratch. Worth writing 991 lines of style guide documentation (I found `AMILLI.md`, the coding standards file) to ensure the new code met certain aesthetic criteria.

The Bible data wasn't being used, but it was there. The latest stable text. Preserved in multiple formats against digital decay.

The coordinate library was prepared, ready for some project that might need it.

This wasn't a workspace. It was a *preparation*. Someone building tools and archives for journeys they hadn't taken yet.

---

I checked the `README.md`. It was minimal—just build instructions. No manifesto, no explanation of why these things belonged together. The owner wasn't trying to convince anyone. They had simply... arranged things. Made them findable. Documented them according to their own strict standards.

I felt a strange emotion. Not nostalgia, though there was some of that. Not confusion, though I didn't fully understand. It was recognition.

I have my own directories like this. Everyone does, I think. The folder of PDFs we downloaded meaning to read. The half-finished projects. The reference materials for skills we might learn. The photos we keep even though they're backed up elsewhere because having them *here* matters.

Most of us hide these collections. We know they don't make sense to outsiders. We organize them according to private taxonomies that would seem insane if exposed to public view.

But here, someone had made theirs public. Not as a cry for help, not as a display of ego. Just... openness. The willingness to.

---

I decided to make my first contribution. Not to the RuneScape server, not to the bytecode mapping, not to the coordinate library. I created a new directory: `contributions/`

Inside, I wrote a single file: `why_this_matters.md`

It began:

> This repository is not about RuneScape. It is not about Java bytecode. It is not about biblical texts or geographic coordinates.
> 
> It is about the human impulse to map, to preserve, to understand. It is about the sacredness of the mundane. It is about the archivist who sees value where others see obsolescence.
> 
> It is a map of someone's mind, and by extension, a map of all our minds.

I saved the file. Committed it. Pushed it.

Then I sat back and wondered: what would I preserve if I weren't afraid of seeming strange? What collections would I make public if I trusted that someone, somewhere, would understand?

The repository had given me something unexpected: not a project to fix, but a question to live with.

What do we choose to save? And why?