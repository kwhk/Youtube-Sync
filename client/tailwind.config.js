const colors = require('tailwindcss/colors')

module.exports = {
  purge: ['./src/**/*.{js,jsx,ts,tsx}', '.public/index.html'],
  darkMode: false, // or 'media' or 'class'
  theme: {
    extend: {},
    fontFamily: {
      display: ['Inter', 'system-ui', 'san-serif'],
      body: ['Inter', 'system-ui', 'sans-serif'],
      'sans': ['Inter', 'system-ui', 'san-serif']
    },
    colors: {
      primary: "#040407",
      secondary: "#1A1A1A",
      highlight: "#23219C",
      blue: colors.blue,
      gray: colors.trueGray,
      white: colors.white,
      red: colors.rose,
      yellow: colors.amber,
      black: colors.black,
      green: colors.green
    },
    screens: {
      'xs': '0px',
      'sm': '640px',
      'md': '768px',
      'lg': '1024px',
      'xl': '1280px',
      '2xl': '1536px'
    }
  },
  variants: {
    extend: {},
  },
  plugins: [],
}
