-- Remove CAC access for abbyoung, suzubara, ralren, kimallen, and garrettqmartin8 using unique sha256 digest of client cert
DELETE FROM client_certs WHERE sha256_digest in (
    '68e7ce8ec4681f4eb1cbebf2f2e3e25eb1edb4675689087b66efabd0b04b330a',
    '87298a2a98a43569d45c359488f9a865aed52a51e60ceed2c5beb69528b05d3e',
    '1bd56ab913eeaf56fc85067cc0d126a400e9513d222551ee8bc1cfb000749934',
    'fcee9e07caf20dbcc2c795652c232d8b554e01e999941f23167d31e80a3ca330',
    '6e987a7f0abbc885a71b470d64be2bd52adfd21a83ea0ad83bd3b4acafa8b93e'
);
