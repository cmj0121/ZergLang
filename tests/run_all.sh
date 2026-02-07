#!/bin/bash
# Run all Zerg verification tests

set -e

ZERG=${ZERG:-./bin/zerg-bootstrap}
TESTS_DIR="tests"

echo "=== Zerg Runtime Verification ==="
echo ""

# Run Go tests first
echo "Step 1: Go Runtime Tests"
cd src/runtime && go test ./... && cd ../..
echo ""

# Run Zerg verification tests
echo "Step 2: Zerg Verification Tests"
for test in verify_strings verify_chars verify_collections verify_classes verify_enums verify_io; do
    echo "Running $test.zg..."
    $ZERG $TESTS_DIR/$test.zg
    echo ""
done

# Run mini lexer
echo "Step 3: Mini Lexer Test (End-to-End)"
$ZERG $TESTS_DIR/mini_lexer.zg
echo ""

echo "=== All Verification Tests Passed! ==="
echo "The Zerg runtime is ready for self-hosting."
