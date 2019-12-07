import pytest
import sys

class CollectTests:
    def __init__(self):
        self.collected = set()

    def pytest_collection_finish(self, session):
        for item in session.items:
            self.collected.add(item.location[0])


collect_tests = CollectTests()
markers = sys.argv[1]
directory = sys.argv[2]
pytest.main(['--collect-only', '-m', markers, '-p', 'no:terminal', directory], plugins=[collect_tests])
for location in collect_tests.collected:
    print(location)
