import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { map, catchError } from 'rxjs/operators';

@Injectable()
export class FileService {

constructor(
    private httpClient: HttpClient,
  ) { }

    postFile(patientId: string, fileToUpload: File): Observable<boolean> {
        const formData: FormData = new FormData();
        formData.append('fileUpload', fileToUpload, fileToUpload.name);
        formData.append('tags', "tag1 tag2");
        //todo add tags
        return this.httpClient
            .post('/api/patient/upload', formData, {params: {"id": patientId}})
            .pipe(map(() => { return true; }));
    }
}