#!/usr/bin/env perl

use strict;
use warnings;
use utf8;
use feature qw/say/;

local $| = 1;
my $s = "X" x 9;
say $s;
sleep(2);
say $s;
sleep(2);
say $s;
sleep(2);

__END__

