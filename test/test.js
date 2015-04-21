(function(has_call, info){
	var fn = this.fn = {};
	if (has_call) {
		putout(info);
	}else{
		puterror();
	}

	function putout(info) {
		console.Write(info);
		return info;
	}

	function puterror() {
		throw "error";
	}

	fn.return_test = function () {
		var test1, test2, test3 = "test";

		test1 = "aster";

		if (fn != null) {
			return test1;
		}

		return test2;
	}
})(false, "error")