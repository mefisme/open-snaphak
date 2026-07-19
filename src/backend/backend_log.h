/* backend_log.h -- the Snapmap+ backend's minimal log sink (OutputDebugStringA + a log file next to
 * the DLL). Independent of the fault-shield's fault_record. OutputDebugStringA -> in-game console +
 * any attached debugger's output sink (the smoke proof surface). */
#ifndef BACKEND_LOG_H
#define BACKEND_LOG_H

#include <windows.h>

/* Set the log path to <dir-of-self>\snapmap-plus\logs\sh_backend.log. Call once from DllMain. */
void backend_set_logpath_from_module(HINSTANCE self);

/* Timestamp + "[snapmap+] " prefix; write to OutputDebugStringA and append to the log file. */
void backend_log(const char *msg);

#endif /* BACKEND_LOG_H */
