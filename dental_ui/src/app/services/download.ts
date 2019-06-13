import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';

@Injectable()
export class FileDownloader {

constructor(
    private http: HttpClient,
  ) { }
  
    downloadFile() {
        return this.http
        .get('/api/download', {responseType: 'blob'})
        .subscribe(response => this.downLoadFile(response /*, "application/ms-excel"*/));
    }

    /**
     * Method is use to download file.
     * @param data - Array Buffer data
     * @param type - type of the document.
     */
    downLoadFile(data: any, type?: string) {
        let blob = new Blob([data] /*, { type: type}*/);
        let blobURL = window.URL.createObjectURL(blob);

        var anchor = document.createElement("a");
        anchor.download = "myfile.png";
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