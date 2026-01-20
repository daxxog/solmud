package supercharger

/*
 ╱╲
╱  ╲
╲  ╱
 ╲╱
*/


/* A Tile is the fundamental building block composing the GameWorld.
 *
 *
 * Tiles consist of:
 *  - 4 byte map identifier.
 *  - 4 equidistant Points (North-West, North-East, South-West, South-East)
 *    NW    NE
 *       ■■   
 *       ■■   
 *    SW    SE
 *  - Points can be connected to any other point in the same Tile.
 *  - Connections do not bear directionality.
 *
 * A Tile is either ACCESSIBLE or INACCESSIBLE. This is determined by:
 *  1. Containing an INACCESSIBLE Interactable makes the Tile INACCESSIBLE.
 *  2. Less than 4 connections makes the Tile ACCESSIBLE, unless (1.) is true.
 *  3. However any connection crossing the middle (SW <-> NE, NW <-> SE) causes the Tile to be INACCESSIBLE.
 *  4. Accessibility is not the same concept as REACHABLE
 *     —an ACCESSIBLE Tile surrounded by INACCESSIBLE tiles is UNREACHABLE (except by Teleportation).
 * 
 * A Tile can contain Interactables.
 * A Tile can contain Players.
 *
 * Tiles have representations in UNICODE:
 *  - (ACCESSIBLE) " "   : Empty
 *  - (ACCESSIBLE) "★"   : Contains Current Player
 *  - (ACCESSIBLE) "☆"   : Contains Other Player
 *  - (ACCESSIBLE) "¤"   : Contains Multiple Interactables
 *  - (ACCESSIBLE) "°"   : Contains Interactable
 *  - (INACCESSIBLE) "△" : Contains INACCESSIBLE Interactable = Tile also INACCESSIBLE
 *  - (INACCESSIBLE) "◬" : Contains Multiple Interactables (at least one is INACCESSIBLE,
                                                          causing the Tile to inherit the INACCESSIBLE state)

┏━┳━┓
┃╳┃ ┃
┗━┻━┛

┏━┓
┃ ┃
┗━┛

┏━┳━━━━━━
┃ ┃     
┣━┻━┓
 ╲ ╱    
  ╳     ┏
 ╱ ╲    ┃
┃   ┣━┳━┫
 ╲ ╱┃ ┃╱
  ━ ┗━┻

┏━┓
┃ ┃
┣━┫
┃ ┃
┗━┛

┏━┳━┓
┃ ┃ ┃
┣━╋━┫
┃ ┃ ┃
┗━┻━┛

┃ ┃
┗━┛

┏━┓
┃ ┃

┏━┓
┃☆┃
┗━┛

┏━┓
┃★┃
┗━┛

┏━┓
┃¤┃
┗━┛

┏━┓
┃°┃
┗━┛

┏━┓
┃△┃
┗━┛

┏━┓
┃◬┃
┗━┛

 */
