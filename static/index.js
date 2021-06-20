const app = new Vue({
	"el": "#app",
	"data": {
		"state": "scan",
		"verifName": "",
		"verifDob": "",
		"verifApproved": false,
		"verifVaccineType": ""
	}
})

const succ = async (decodedText, decodedResult) => {
	if (app.state !== "scan") {
		return
	}

	const r = /shc:\/([0-9]+)/.exec(decodedText)

	if (r == null) {
		app.state = "error"
		app.errorMessage = "Invalid code. This is not a vaccination record QR code."
		return
	}

	const res = await fetch("/check", {
		"method": "POST",
		"cache": "no-cache",
		"headers": {
			"Content-Type": "application/json"
		},
		"body": r[1]
	})

	if (res.status != 200) {
		const txt = await res.text()

		app.state = "error"
		app.errorMessage = txt

		return
	}

	const parsed = await res.json()

	console.log(parsed)

	app.verifApproved = parsed.approved
	app.verifName = parsed.fullName
	app.verifDob = parsed.patientDob
	app.verifVaccineType = parsed.vaccineType
	app.state = "ok"
}

const fail = (err) => {
	//alert("Error: " + err)
	if (err.toString().includes("No MultiFormat Readers")) {
		return
	}

	console.error(err)
}

const scanner = new Html5QrcodeScanner("reader", {
	"fps": 10,
	"qrbox": 250,
}, false)

scanner.render(succ, fail)
