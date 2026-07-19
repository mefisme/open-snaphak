/* snapstack.h -- the BACKEND-hosted SnapStack implementation: the stores (stack-of-stacks +
 * named groups) and the 20 `sh <subcommand>` console handlers, in pure C.
 *
 * This module is the SOLE SnapStack implementation. It registers the 20 command names on the shared
 * sh_iface cmd-map (which overwrites on duplicate name -- see snapmap_plus_iface.c's iface_register_cmd)
 * from ui_bridge.c right after the interface object is created, before the frontend loads -- so every
 * `sh <subcommand>` runs this module's handlers. The frontend never registers a copy of its own.
 *
 * Stores are file-static singletons here (one stack-of-stacks + one group map), compiled into
 * XINPUT1_3.dll alongside everything else that touches them.
 *
 * Clean-room: ported from our own RE; the JSON structural work lives in json_patch.c. Zero OG
 * SnapHak bytes.
 */
#ifndef BACKEND_SNAPSTACK_H
#define BACKEND_SNAPSTACK_H

#include "snapmap_plus_iface.h"

/* Register the 20 SnapStack subcommands PLUS `snapstack_diag` (a backend-exclusive diagnostic -- see
 * snapstack.c) on the interface's cmd-map (+0x188). Safe to call once, early (before the frontend loads)
 * -- ctx = iface for every handler. */
void sh_register_snapstack_commands_backend(sh_iface *iface);

/* The Entities-tab "Push to stack 0" context-menu action, exposed for the push_to_stack vtable slot
 * (iface_engine.c) so ANY frontend (not just one compiled in-process with this module) can push onto
 * the SAME backend-owned stack a `sh psel`/`sh popsel`/etc. console command run afterward will see.
 * Dedup-on-push, same as every other path into the stores. */
void sh_snapstack_push_ids_backend(int index, const int *ids, int count);

/* The Entities-tab "Clear stack 0" context-menu action, exposed for the clear_stack vtable slot
 * (iface_engine.c) so ANY frontend can empty the SAME backend-owned stack a `sh cstk` console command
 * would -- lets a user work purely from the UI without ever touching the DOOM console. Returns the number
 * of ids that were on the stack before clearing (so the caller can toast a count, same as `sh cstk` itself). */
int sh_snapstack_clear_stack_backend(int index);

#endif /* BACKEND_SNAPSTACK_H */
