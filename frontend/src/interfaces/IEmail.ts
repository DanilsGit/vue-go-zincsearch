interface IEmail {
    "@timestamp": string;
    content: string;
    content_transfer_encoding: string;
    content_type: string;
    date: string;
    from: string;
    message_id: string;
    mime_version: string;
    subject: string;
    to: string;
    x_bcc: string;
    x_cc: string;
    x_filename: string;
    x_folder: string;
    x_from: string;
    x_origin: string;
    x_to: string;
    subject_mark: string;
    content_mark: string;
}

export interface IHit {
    "@timestamp": string;
    _id: string;
    _index: string;
    _score: number;
    _source: IEmail;
}
