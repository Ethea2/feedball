# Feedball

A 1v1 competitive action game where two players battle it out in a frog terrarium, trying to feed the opposing player's frog a poisonous fly — before their own frog gets fed first.

---

## Overview

Feedball is a fast-paced, physics-driven 1v1 game built around one ball, two frogs, and a whole lot of ricochet. Players take turns as the **Attacker** and **Defender**, shooting a ball across a multi-level arena, working angles off walls to open and score into the opposing frog's mouth.

The first player to score 5 goals wins.

---

## Lore

At Shaw Tower unit 1507, an evil wizard shrinks two roommates to insect size and traps them in a frog terrarium. The only way back to normal size? Defeat each other's frog.

---

## Gameplay

### Roles
At any given moment, one player is the **Attacker** (ball in hand) and the other is the **Defender**. Roles shift dynamically as possession changes.

### The Arena
A flat-bottomed plane with platforms of varying heights. Walls are everywhere — and every wall is a potential angle.

---

## The Ball

The ball travels in a straight line, unaffected by gravity, and ricochets off walls at right angles. It slows to a capped speed after 3 bounces.

### Ball States

| State | Description |
|---|---|
| **INPLAY** | Held by the Attacker. Bounces off the Defender like a wall. Becomes **LIVE** after 2 bounces. |
| **LIVE** | No owner. First player to touch it gains possession and becomes the Attacker. |
| **GROSS** | Spit out by the frog. Cannot be grabbed. Knocks back and briefly slows any player it hits. Becomes **LIVE** after its first wall bounce. |

---

## Scoring & The Frog

The **goal** is the frog's mouth — about 4–5 ball widths wide.

- The mouth is naturally **closed** and acts as a wall.
- Any ball contact opens the mouth for **10 seconds**.
- Score by shooting the ball into the open mouth → **–1 HP** to the frog.
- The frog has **5 HP**. First player to bring the opposing frog to 0 wins.

### Frog Behavior
- If the Attacker gets **too close** to an open frog mouth, the frog gets scared and **snaps shut**. Get creative with your angles.
- On a successful goal (except the final one), the frog **eats** the ball and **spits it out** as a **GROSS** ball — aimed directly at the Attacker who just scored.

---

## Controls & Mechanics

### Movement
- **Run** to move around the arena.
- **Jump** to block INPLAY balls or grab LIVE balls mid-air.
- Movement speed and jump height are **slightly reduced** when holding the ball.

### Shooting
The Attacker can shoot in 6 directions:

```
↖  ↑  ↗
←     →
↙  ↓  ↘
```

### Floating
When the Attacker shoots, there's a brief input window during which they **float** before falling normally. Use this to fine-tune your shot direction.

### Fast Fall
The **Defender** can press **↓** to fast fall — dropping faster than normal to reposition quickly.

---

## Game Start — Jump Ball

At the start of each game:

1. The ball spawns **LIVE** in the center of the arena.
2. Both players are locked in place on opposite sides.
3. After a **random 3–6 second** delay, the ball jumps upward.
4. Movement unlocks — the first player to reach and grab the ball becomes the Attacker.

---

## Sample Round Loop

```
Grab the jump ball
→ Maneuver and shoot at the opposing frog
→ Frog mouth opens
→ Shoot the ball into the open mouth (–1 HP)
→ Frog spits out a GROSS ball — dodge it
→ Compete for possession
→ Repeat until the opposing frog hits 0 HP
```

---

## Win Condition

Deplete the opposing frog's health bar from **5 to 0**. First to 5 goals wins the game.
