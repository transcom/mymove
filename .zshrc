# Alias for clean milmove 
alias mmclean='make clean server_generate db_dev_reset db_dev_e2e_populate'

#hopefully stops nix from getting broken by updates to macos
if [ -e '/nix/var/nix/profiles/default/etc/profile.d/nix-daemon.sh' ]; then
  . '/nix/var/nix/profiles/default/etc/profile.d/nix-daemon.sh'
fi
# End Nix fix
