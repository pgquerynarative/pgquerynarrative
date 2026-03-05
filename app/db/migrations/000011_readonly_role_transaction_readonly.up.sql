-- Enforce read-only at session level for the readonly role.
ALTER ROLE pgquerynarrative_readonly SET default_transaction_read_only = on;
