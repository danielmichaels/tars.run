module.exports = {
  content: ["./ui/**/*tmpl"],
  theme: {
    extend: {},
  },
  plugins: [
      require('@tailwindcss/aspect-ratio'),
      require('@tailwindcss/forms'),
  ],
}
