from unittest.mock import mock_open
from unittest.mock import patch

from zergb.compiler import ZergBootstrap


class TestBootstrap:
    def setup_method(self):
        self.bootstrap = ZergBootstrap()

    def test_nop(self):
        src = 'fn main() { nop }'

        with patch('builtins.open', mock_open(read_data=src)):
            self.bootstrap._source_to_ir('fake')
