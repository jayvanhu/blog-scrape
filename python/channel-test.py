from unittest import TestCase
import unittest

from channel import Channel

class TestChannel(TestCase):

	def setUp(self):
		self.chan = Channel()

	def test_inits_unfinished(self):
		self.assertFalse(self.chan.finished())

	def test_unfinished_after_put(self):
		self.chan.put(1)
		self.assertFalse(self.chan.finished())

	def test_unfinished_after_put_and_get(self):
		self.chan.put(1)
		self.chan.get()
		self.assertFalse(self.chan.finished())

	def test_unfinished_after_put_and_close(self):
		self.chan.put(1)
		self.chan.close()
		self.assertFalse(self.chan.finished())

	def test_finished_after_close(self):
		self.chan.close()
		self.assertTrue(self.chan.finished())

	def test_finished_after_close_and_get(self):
		self.chan.put(1)
		self.chan.get()
		self.chan.close()
		self.assertTrue(self.chan.finished())

if __name__ == '__main__':
	unittest.main()
