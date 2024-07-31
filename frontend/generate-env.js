import { writeFileSync } from 'fs';

const envVariables = [
  'API_URL'
];

let envContent = '';

envVariables.forEach(variable => {
    console.log(process.env[variable]);
    if (process.env[variable]) {
        envContent += `VITE_${variable}=${process.env[variable]}\n`;
        envContent += `found_${variable}=${process.env[variable]}\n`;
    }else{
        envContent += `not_found\n`;
    }
});

writeFileSync('.env', envContent);
