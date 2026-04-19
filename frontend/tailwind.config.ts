import type { Config } from "tailwindcss";

const config: Config = {
  content: ["./app/**/*.{ts,tsx}", "./components/**/*.{ts,tsx}"],
  theme: {
    extend: {
      colors: {
        navy: "#0B1A33",
        gold: "#C7A64B"
      },
      boxShadow: {
        midas: "0 10px 25px rgba(11,26,51,0.08)"
      }
    }
  },
  plugins: []
};

export default config;
