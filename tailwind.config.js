/** @type {import('tailwindcss').Config} */
const plugin = require("tailwindcss/plugin");

module.exports = {
    content: ["./internal/render/**/*.templ", "./internal/render/styles/*.go"],
    theme: {
        colors: {
            white: "#ffffff",
            primary: {
                100: "#633103",
                200: "#954a04",
                300: "#c66206",
                400: "#df6e06",
                500: "#f98820",
                600: "#faa251",
                700: "#fbb06a",
                800: "#fbbd83",
                900: "#fdddbf",
            },
            secondary: {
                100: "#5c0a0a",
                200: "#a32929",
                300: "#d14747",
                400: "#db5757",
                500: "#f28787",
                600: "#f49a9a",
                700: "#f5a8a8",
                800: "#f9c8c8",
                900: "#fbdada",
            },
            grey: {
                100: "#34302d",
                200: "#3c3834",
                300: "#70655c",
                400: "#9c958c",
                500: "#b8b2ad",
                600: "#cfccc9",
                700: "#e0dedc",
                800: "#f8f7f7",
                900: "#fdfcfc",
            },
            red: {
                100: "#6d0303",
                200: "#aa0909",
                300: "#c11515",
                400: "#d61f1f",
                500: "#e83030",
                600: "#fb8383",
                700: "#fdb5b5",
                800: "#fed7d7",
                900: "#fee6e6",
            },
            green: {
                200: "#0b750b",
                500: "#70c270",
                800: "#dffcdf",
            },
            transparent: {},
        },
        fontSize: {
            sm: "0.750rem",
            base: "1rem",
            xl: "1.333rem",
            "2xl": "1.777rem",
            "3xl": "2.369rem",
            "4xl": "3.158rem",
            "5xl": "4.210rem",
        },
        fontFamily: {
            heading: "Inter",
            body: "Inter",
        },
        fontWeight: {
            normal: "400",
            bold: "700",
        },
        extend: {},
    },
    plugins: [
        plugin(function ({ addVariant }) {
            addVariant("hocus", ["&:hover", "&:focus"]);
        }),
    ],
};
