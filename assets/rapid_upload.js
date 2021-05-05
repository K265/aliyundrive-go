function i(n, t) {
  var e = (65535 & n) + (65535 & t);
  return (n >> 16) + (t >> 16) + (e >> 16) << 16 | 65535 & e
}
function o(n, t, e, a, r, o) {
  return i((l = i(i(t, n), i(a, o))) << (A = r) | l >>> 32 - A, e);
  var l, A
}
function l(n, t, e, a, r, i, l) {
  return o(t & e | ~t & a, n, t, r, i, l)
}
function A(n, t, e, a, r, i, l) {
  return o(t & a | e & ~a, n, t, r, i, l)
}
function s(n, t, e, a, r, i, l) {
  return o(t ^ e ^ a, n, t, r, i, l)
}
function c(n, t, e, a, r, i, l) {
  return o(e ^ (t | ~a), n, t, r, i, l)
}

function p(n, t) {
  var e, a, r, o, p;
  n[t >> 5] |= 128 << t % 32,
  n[14 + (t + 64 >>> 9 << 4)] = t;
  var u = 1732584193
    , d = -271733879
    , f = -1732584194
    , m = 271733878;
  for (e = 0; e < n.length; e += 16)
      a = u,
      r = d,
      o = f,
      p = m,
      u = l(u, d, f, m, n[e], 7, -680876936),
      m = l(m, u, d, f, n[e + 1], 12, -389564586),
      f = l(f, m, u, d, n[e + 2], 17, 606105819),
      d = l(d, f, m, u, n[e + 3], 22, -1044525330),
      u = l(u, d, f, m, n[e + 4], 7, -176418897),
      m = l(m, u, d, f, n[e + 5], 12, 1200080426),
      f = l(f, m, u, d, n[e + 6], 17, -1473231341),
      d = l(d, f, m, u, n[e + 7], 22, -45705983),
      u = l(u, d, f, m, n[e + 8], 7, 1770035416),
      m = l(m, u, d, f, n[e + 9], 12, -1958414417),
      f = l(f, m, u, d, n[e + 10], 17, -42063),
      d = l(d, f, m, u, n[e + 11], 22, -1990404162),
      u = l(u, d, f, m, n[e + 12], 7, 1804603682),
      m = l(m, u, d, f, n[e + 13], 12, -40341101),
      f = l(f, m, u, d, n[e + 14], 17, -1502002290),
      u = A(u, d = l(d, f, m, u, n[e + 15], 22, 1236535329), f, m, n[e + 1], 5, -165796510),
      m = A(m, u, d, f, n[e + 6], 9, -1069501632),
      f = A(f, m, u, d, n[e + 11], 14, 643717713),
      d = A(d, f, m, u, n[e], 20, -373897302),
      u = A(u, d, f, m, n[e + 5], 5, -701558691),
      m = A(m, u, d, f, n[e + 10], 9, 38016083),
      f = A(f, m, u, d, n[e + 15], 14, -660478335),
      d = A(d, f, m, u, n[e + 4], 20, -405537848),
      u = A(u, d, f, m, n[e + 9], 5, 568446438),
      m = A(m, u, d, f, n[e + 14], 9, -1019803690),
      f = A(f, m, u, d, n[e + 3], 14, -187363961),
      d = A(d, f, m, u, n[e + 8], 20, 1163531501),
      u = A(u, d, f, m, n[e + 13], 5, -1444681467),
      m = A(m, u, d, f, n[e + 2], 9, -51403784),
      f = A(f, m, u, d, n[e + 7], 14, 1735328473),
      u = s(u, d = A(d, f, m, u, n[e + 12], 20, -1926607734), f, m, n[e + 5], 4, -378558),
      m = s(m, u, d, f, n[e + 8], 11, -2022574463),
      f = s(f, m, u, d, n[e + 11], 16, 1839030562),
      d = s(d, f, m, u, n[e + 14], 23, -35309556),
      u = s(u, d, f, m, n[e + 1], 4, -1530992060),
      m = s(m, u, d, f, n[e + 4], 11, 1272893353),
      f = s(f, m, u, d, n[e + 7], 16, -155497632),
      d = s(d, f, m, u, n[e + 10], 23, -1094730640),
      u = s(u, d, f, m, n[e + 13], 4, 681279174),
      m = s(m, u, d, f, n[e], 11, -358537222),
      f = s(f, m, u, d, n[e + 3], 16, -722521979),
      d = s(d, f, m, u, n[e + 6], 23, 76029189),
      u = s(u, d, f, m, n[e + 9], 4, -640364487),
      m = s(m, u, d, f, n[e + 12], 11, -421815835),
      f = s(f, m, u, d, n[e + 15], 16, 530742520),
      u = c(u, d = s(d, f, m, u, n[e + 2], 23, -995338651), f, m, n[e], 6, -198630844),
      m = c(m, u, d, f, n[e + 7], 10, 1126891415),
      f = c(f, m, u, d, n[e + 14], 15, -1416354905),
      d = c(d, f, m, u, n[e + 5], 21, -57434055),
      u = c(u, d, f, m, n[e + 12], 6, 1700485571),
      m = c(m, u, d, f, n[e + 3], 10, -1894986606),
      f = c(f, m, u, d, n[e + 10], 15, -1051523),
      d = c(d, f, m, u, n[e + 1], 21, -2054922799),
      u = c(u, d, f, m, n[e + 8], 6, 1873313359),
      m = c(m, u, d, f, n[e + 15], 10, -30611744),
      f = c(f, m, u, d, n[e + 6], 15, -1560198380),
      d = c(d, f, m, u, n[e + 13], 21, 1309151649),
      u = c(u, d, f, m, n[e + 4], 6, -145523070),
      m = c(m, u, d, f, n[e + 11], 10, -1120210379),
      f = c(f, m, u, d, n[e + 2], 15, 718787259),
      d = c(d, f, m, u, n[e + 9], 21, -343485551),
      u = i(u, a),
      d = i(d, r),
      f = i(f, o),
      m = i(m, p);
  return [u, d, f, m]
}

function d(n) {
  var t, e = [];
  for (e[(n.length >> 2) - 1] = void 0,
  t = 0; t < e.length; t += 1)
      e[t] = 0;
  var a = 8 * n.length;
  for (t = 0; t < a; t += 8)
      e[t >> 5] |= (255 & n.charCodeAt(t / 8)) << t % 32;
  return e
}

function u(n) {
  var t, e = "", a = 32 * n.length;
  for (t = 0; t < a; t += 8)
      e += String.fromCharCode(n[t >> 5] >>> t % 32 & 255);
  return e
}

function h(n) {
  return u(p(d(n), 8 * n.length))
}

function f(n) {
  var t, e, a = "0123456789abcdef", r = "";
  for (e = 0; e < n.length; e += 1)
      t = n.charCodeAt(e),
      r += a.charAt(t >>> 4 & 15) + a.charAt(15 & t);
  return r
}

function g(n) {
  return f(h(n))
}

function calcProof(access_token, file_size) {
  var e = access_token // "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9..."
  , a = BigInt("0x".concat(g(e).slice(0, 16))) // 0xaa4b5fde957fd92b
  , r = BigInt(file_size) // 1477419708
  , i = r ? a%r : BigInt(0) // 915335463
  // , o = t.file.slice(i.toNumber(), Math.min(i.plus(8).toNumber(), t.file.size)) // [915335463, 915335471)
  // return l.readAsDataURL(o),
  console.log({a,r,i})
}

