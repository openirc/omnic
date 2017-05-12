==============
Omnic protocol
==============

Why? We need a way for Omnic to communicate with the browser, IRC, and the user.

For these examples, we will assume Ember as the client, though it can be implemented in basically anything that supports websockets.


- O = Command should be parsed by Omnic
- E = Command should be parsed by Ember
- I = Command should be sent to IRC


Ember -> Omnic commands
=======================

Authenticating:

:COMMAND: O A <token>
:RETURNS: E A <int64>


RETURN VALUE MEANINGS:

- -1 = API Failure
- 0  = Access denied
- >0 = User ID / Success

Omnic -> Ember commands
=======================
:COMMAND: E A <int64>
