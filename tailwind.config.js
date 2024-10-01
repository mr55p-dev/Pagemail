module.exports = {
    content: ["render/*.templ"],
    prefix: "tw-",
    theme: {
        extend: {
            colors: {
                brand: {
                    '50': '#eefaff',
                    '100': '#d9f4ff',
                    '200': '#bbecff',
                    '300': '#8ce3ff',
                    '400': '#56cfff',
                    '500': '#2fb4ff',
                    '600': '#1898f8',
                    '700': '#117fe5',
                    '800': '#1565b8',
                    '900': '#175691',
                    '950': '#133558',
                },
                background: {
                    primary: '#ffffff',
                    secondary: '#f5f5f5'
                }
            }
        },
    },
    plugins: [],
};
