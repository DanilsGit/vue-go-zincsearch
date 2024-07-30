import { writeFileSync } from 'fs';

const envVariables = [
  'API_URL'
];

let envContent = '';

envVariables.forEach(variable => {
    console.log(process.env[variable]);
    envContent += `VITE_${variable}=${process.env[variable]}\n`;
});

writeFileSync('.env', envContent);
