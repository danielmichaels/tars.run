module.exports = {
  content: ["./assets/**/*tmpl", "./assets/**/*templ"],
  theme: {
    extend: {},
  },
  plugins: [
      require('@tailwindcss/aspect-ratio'),
      require('@tailwindcss/forms'),
  ],
}
