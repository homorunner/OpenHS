# Hearthstone Tags

This document tracks the implementation status of various Hearthstone tags in the OpenHS project.

## Implemented Tags

The following tags have been implemented and tested:

| Tag | Description | Implemented | Tested |
|-----|-------------|-------------|--------|
| TAG_POISONOUS | Destroys any minion damaged by this | ✅ | ✅ |
| TAG_WINDFURY | Can attack twice per turn | ✅ | ✅ |
| TAG_CHARGE | Can attack on the turn it's played | ✅ | ✅ |
| TAG_RUSH | Can attack minions on the turn it's played | ✅ | ✅ |
| TAG_LIFESTEAL | Damage dealt heals your hero | ✅ | ✅ |
| TAG_FROZEN | Miss next possible attack | ✅ | ✅ |


## Unimplemented Tags

The following tags are defined but not fully implemented or tested:

| Tag | Description | Implementation Notes |
|-----|-------------|---------------------|
| TAG_TAUNT | Forces enemies to attack this minion | Defined but not implemented |
| TAG_DIVINE_SHIELD | Absorbs the next damage instance | Defined but not implemented |
| TAG_STEALTH | Cannot be targeted by opponents | Defined but not implemented |
| TAG_DEATHRATTLE | Triggers an effect when destroyed | Defined but not implemented |
| TAG_BATTLECRY | Triggers an effect when played | Defined but not implemented |
| TAG_REBORN | Returns to life with 1 Health | Defined but not implemented |
| TAG_DORMANT | Cannot act for a set number of turns | Defined but not implemented |
| TAG_SPELLPOWER | Increases spell damage | Defined but not implemented |
| TAG_CANT_ATTACK | Cannot attack at all | Defined but not implemented |
| TAG_CANT_BE_TARGETED | Cannot be targeted by spells/hero powers | Defined but not implemented |
| TAG_IMMUNE | Cannot take damage | Defined but not implemented |

## Development Guidelines

When implementing a new tag:
1. Add the tag constant to `internal/game/tag.go`
2. Implement the game logic to handle the tag
3. Add tests for the tag in the tests directory
4. Update this document to mark the tag as implemented and tested 