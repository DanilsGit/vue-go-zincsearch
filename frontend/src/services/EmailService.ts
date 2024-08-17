import axios from 'axios';
import { ref, Ref } from 'vue';
import {IHit} from '../interfaces/IEmail';

export class EmailService {
  private emails:Ref<Array<IHit>>;
  private totalEmails:Ref<number> = ref(0);

  constructor() {
    this.emails = ref<Array<IHit>>([]);
  }

  getEmails():Ref<Array<IHit>> {
    return this.emails;
  }

  getTotalEmails():Ref<number> {
    return this.totalEmails;
  }

  async fetchAllEmails(sort:string, offset:number, limit:number):Promise<void> {
    try {
        const url = import.meta.env.VITE_API_URL + '/getAll';
        console.log("dev propouse", url);
        
        const response = await axios.get(url, {
          params: {
            from: offset.toString(),
            max: limit.toString(),
            sort : sort
          }
        });
        
        response.data.hits.hits.forEach((email:IHit) => {
          email._source.date = this.parseDate(email._source.date);
          email._source.x_from = this.parseName(email._source.x_from);
        });
        this.emails.value = response.data.hits.hits;
        this.totalEmails.value = response.data.hits.total.value;
    } catch (error) {
        console.log("Error fetching emails");
        console.log(error);
    }
  }

  async fetchSearchEmails(search:string, sort:string, offset:number, limit:number):Promise<void> {
    try {
        const url = import.meta.env.VITE_API_URL + '/search';
        const response = await axios.get(url, {
          params: {
            type: 'matchphrase',
            from: offset.toString(),
            max: limit.toString(),
            sort : sort,
            search : search
          }
        });

        if (!response.data.hits) response.data.hits = {hits: [], total: {value: 0}};
        
        response.data.hits.hits.forEach((email:IHit) => {
          console.log(email._source.date);
        });

        response.data.hits.hits.forEach((email:IHit) => {
          email._source.date = this.parseDate(email._source.date);
          email._source.x_from = this.parseName(email._source.x_from);
        });
        this.emails.value = response.data.hits.hits;
        this.totalEmails.value = response.data.hits.total.value;
    } catch (error) {
        // console.log("Error fetching emails");
        // console.log(error);
    }
  }

  parseDate(date:string):string {
    const dateSplit = new Date(date).toLocaleString().split(',')[0].split('/');
    const day = Number(dateSplit[0]) < 10 ? '0' + dateSplit[0] : dateSplit[0];
    const month = Number(dateSplit[1]) < 10 ? '0' + dateSplit[1] : dateSplit[1];
    const year = Number(dateSplit[2]) < 1000 ? '000' + dateSplit[2] : dateSplit[2];
    return `${day}/${month}/${year}`;
  }

  parseName(name:string):string {
    name = name.split('<')[0];
    if (name.split( '@' ).length > 1) name = name.split( '@' )[0];
    if (name.split ( ' ' ).length > 1) name = name.split ( ' ' )[0] + ' ' + name.split ( ' ' )[1];
    if (name.length === 0) name = 'Unknown';
    return name;
  }


}