package harvesterd

import (
	"harvesterd/intf"
)

import . "launchpad.net/gocheck"

type WriterGroupSuite struct {
	wg *WriterGroup
	wA *MockWriter
	wB *MockWriter
}

var _ = Suite(&WriterGroupSuite{})

func (self *WriterGroupSuite) SetUpTest(c *C) {
	self.wA = new(MockWriter)
	self.wB = new(MockWriter)

	self.wg = NewWriterGroup()
	self.wg.SetWriters([]intf.Writer{self.wA, self.wB})
}

func (self *WriterGroupSuite) TestSetup(c *C) {
	self.wg.Setup()
	c.Assert(self.wA.setupCount, Equals, 1)
	c.Assert(self.wB.setupCount, Equals, 1)
}

func (self *WriterGroupSuite) TestBoot(c *C) {
	self.wg.Boot()
	c.Assert(self.wA.bootCount, Equals, 1)
	c.Assert(self.wB.bootCount, Equals, 1)
}

func (self *WriterGroupSuite) TestResetCounters(c *C) {
	self.wg.ResetCounters()
	c.Assert(self.wA.resetCountersCount, Equals, 1)
	c.Assert(self.wB.resetCountersCount, Equals, 1)
}

func (self *WriterGroupSuite) TestTeardown(c *C) {
	self.wg.Teardown()
	c.Assert(self.wA.teardownCount, Equals, 1)
	c.Assert(self.wB.teardownCount, Equals, 1)
}

func (self *WriterGroupSuite) TestGetCounters(c *C) {
	o, f, t, h := self.wg.GetCounters()

	c.Assert(o, Equals, int32(2))
	c.Assert(f, Equals, int32(4))
	c.Assert(t, Equals, int32(6))
	c.Assert(h, Equals, int32(8))
}

func (self *WriterGroupSuite) TestIsAliveBothAlive(c *C) {
	self.wA.isAlive = true
	self.wB.isAlive = true

	c.Assert(self.wg.IsAlive(), Equals, true)

	c.Assert(self.wA.isAliveCount, Equals, 1)
	c.Assert(self.wB.isAliveCount, Equals, 0)
}

func (self *WriterGroupSuite) TestIsAliveOneAlive(c *C) {
	self.wA.isAlive = false
	self.wB.isAlive = true

	c.Assert(self.wg.IsAlive(), Equals, true)

	c.Assert(self.wA.isAliveCount, Equals, 1)
	c.Assert(self.wB.isAliveCount, Equals, 1)
}

func (self *WriterGroupSuite) TestIsAliveBothDeath(c *C) {
	self.wA.isAlive = false
	self.wB.isAlive = false

	c.Assert(self.wg.IsAlive(), Equals, false)

	c.Assert(self.wA.isAliveCount, Equals, 1)
	c.Assert(self.wB.isAliveCount, Equals, 1)
}

type MockWriter struct {
	setupCount         int
	bootCount          int
	resetCountersCount int
	teardownCount      int
	isAlive            bool
	isAliveCount       int
}

func (self *MockWriter) IsAlive() bool {
	self.isAliveCount++

	return self.isAlive
}

func (self *MockWriter) Setup() {
	self.setupCount++
}

func (self *MockWriter) Boot() {
	self.bootCount++
}

func (self *MockWriter) GetCounters() (int32, int32, int32, int32) {
	return 1, 2, 3, 4
}

func (self *MockWriter) ResetCounters() {
	self.resetCountersCount++
}

func (self *MockWriter) Teardown() {
	self.teardownCount++
}
