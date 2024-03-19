#! /usr/bin/env python3
import argparse
import os

from zergb.compiler import ZergBootstrap


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument('source', help='source file to compile')
    parser.add_argument('-o', '--output', help='output file')
    args = parser.parse_args()

    bootstrap = ZergBootstrap()

    with open(args.source) as fd:
        src = fd.read()

        if args.output is None:
            path = os.path.basename(args.source)
            name, _ = os.path.splitext(path)
            args.output = f'{name}.o'

        with open(args.output, 'wb') as f:
            obj = bootstrap.compile(src)
            f.write(obj)


if __name__ == '__main__':
    main()
