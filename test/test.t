use strict;
use warnings;
use utf8;
use FindBin;
use Test::More;

my $result = system("go build $FindBin::Bin/../cmd/floodgate/floodgate.go");
die "Failed to build" if $result != 0;

subtest 'Missing mandatory arguments' => sub {
    ok system("$FindBin::Bin/floodgate") != 0;
};

subtest 'with byte threshold' => sub {
    my $s = '';
    for my $i (1..12) {
        $s .= sprintf("%s\n", "X" x 99); # 100 bytes line
    } # 1200 bytes

    subtest 'should flush' => sub {
        my $got = `echo \"$s\" | $FindBin::Bin/floodgate --threshold=1024`;
        is length $got, 1100;
    };

    subtest 'should not flush' => sub {
        my $got = `echo \"$s\" | $FindBin::Bin/floodgate --threshold=65536`;
        is length $got, 0;
    };
};

subtest 'with interval' => sub {
    subtest 'should flush' => sub {
        my $got = `perl $FindBin::Bin/canaria.pl | $FindBin::Bin/floodgate --interval=1`;
        is length $got, 30;
    };

    subtest 'should not flush' => sub {
        my $got = `perl $FindBin::Bin/canaria.pl | $FindBin::Bin/floodgate --interval=10`;
        is length $got, 0;
    }
};

subtest 'stderr' => sub {
    subtest 'should flush' => sub {
        my $got = `echo \"foo\" | $FindBin::Bin/floodgate --threshold=1 2>/dev/null`;
        is length $got, 4;
    };

    subtest 'should not flush' => sub {
        my $got = `echo \"foo\" | $FindBin::Bin/floodgate --threshold=1 --stderr 2>/dev/null`;
        is length $got, 0;
    };
};

subtest 'change concatination char' => sub {
    my $got = `echo \"foo\nbar\" | $FindBin::Bin/floodgate --threshold=1 -c='Z'`;
    is $got, 'fooZbarZ';
};

done_testing;

