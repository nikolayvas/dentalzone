import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { map, catchError } from 'rxjs/operators';

@Injectable()
export class FileService {

constructor(
    private httpClient: HttpClient,
  ) { }

    postFile(patientId: string, fileToUpload: File, tags?: string[]): Observable<boolean> {
        const formData: FormData = new FormData();
        formData.append('fileUpload', fileToUpload, fileToUpload.name);
        formData.append('tags', tags.join(" "));
        //todo add tags
        return this.httpClient
            .post('/api/patient/upload', formData, {params: { "id": patientId }})
            .pipe(map(() => { return true; }));
    }

    getFilesByTags(patientId: string, tags: string[]): Observable<IFileModel[]>  {
        return this.httpClient
            .post('/api/patient/files', tags, {params: { "id": patientId }})
            .pipe(map(n => { return n as IFileModel[]; }));
    }

    downloadFile(patientId: string, file: IFileModel ) {
        return this.httpClient
        .post('/api/patient/download', file.id, {responseType: 'blob', params: { "id": patientId }})
        .subscribe(response => this.doDownload(response, file.name /*, "application/ms-excel"*/));
    }

    /**
     * Method is use to download file.
     * @param data - Array Buffer data
     * @param type - type of the document.
     */
    private doDownload(data: any, fileName: string, type?: string) {
        let blob = new Blob([data] /*, { type: type}*/);
        let blobURL = window.URL.createObjectURL(blob);

        var anchor = document.createElement("a");
        anchor.download = fileName;
        anchor.href = blobURL;
        anchor.click();

        /*
        let pwa = window.open(url);
        if (!pwa || pwa.closed || typeof pwa.closed == 'undefined') {
            alert( 'Please disable your Pop-up blocker and try again.');
        }
        */
    }
}

export interface IFileModel {
    id: string,
    name: string,
    size: string,
};