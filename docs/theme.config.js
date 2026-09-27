/**
 * @type {import("nextra-theme-docs").DocsThemeConfig}
 */
export default {
  project: { link: "https://github.com/bongani-m/hardhatkv" },
  docsRepositoryBase: "https://github.com/bongani-m/hardhatkv/tree/master/docs",
  useNextSeoProps() {
    return { titleTemplate: "%s – HardhatKV" };
  },
  primaryHue: { dark: 38, light: 36 },
  logo: (
    <>
      <span className="font-bold" style={{ marginRight: 8 }}>
        HardhatKV
      </span>
      <span className="text-gray-600 font-normal hidden md:inline">
        A Redis-compatible key-value server
      </span>
    </>
  ),
  head: (
    <>
      <meta name="msapplication-TileColor" content="#ffffff" />
      <meta name="theme-color" content="#ffffff" />
      <meta name="viewport" content="width=device-width, initial-scale=1.0" />
      <meta httpEquiv="Content-Language" content="en" />
      <meta
        name="description"
        content="HardhatKV: a Redis-compatible key-value server written in Go"
      />
      <meta
        name="og:description"
        content="HardhatKV: a Redis-compatible key-value server written in Go"
      />
      <meta name="twitter:card" content="summary_large_image" />
      <meta name="twitter:site:domain" content="https://github.com/bongani-m/hardhatkv" />
      <meta name="twitter:url" content="https://github.com/bongani-m/hardhatkv" />
      <meta
        name="og:title"
        content="HardhatKV: a Redis-compatible key-value server written in Go"
      />
      <meta name="apple-mobile-web-app-title" content="HardhatKV" />
    </>
  ),
  navigation: true,
  footer: { text: <>GPL-3.0 {new Date().getFullYear()} © HardhatKV.</> },
  editLink: { text: "Edit this page on GitHub" },
  unstable_faviconGlyph: "⛑",
};
