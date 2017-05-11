import { parse } from "babylon"

let code = ""
process.stdin.on("data", (d) => {
  code += d
}).on("end", () => {
  console.log(JSON.stringify(parse(code)))
}).setEncoding("utf8")
