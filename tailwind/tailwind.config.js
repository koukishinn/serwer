/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
	'../www/**/*.{html,tmpl}'
  ],
  theme: {
      extend: {
		  maxHeight: {
			  '112': '28rem',
			  '128': '32rem',
			  '144': '36rem',
			  '160': '40rem',
		  }
	  },
  },
	plugins: [],
}

