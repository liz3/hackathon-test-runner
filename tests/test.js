describe("some random suite", () => {
	var a;
	it("this is a case", () => {
		a = true;
		expect(a).toBe(true)
	})
	it("another test with a function", function() {
		expect("not").toBe("not")
	})
	it("arrow test case", () => expect(1).toBe(2-1))
})

const a = () => {
	describe("another suite", () => {
		it("test", () => {

		})
	})
}