import { Component } from '@angular/core';
import { FileService } from '../../services/file.service';
import { ActivatedRoute } from '@angular/router';
import { take } from 'rxjs/operators';

@Component({
    selector: 'upload',
    templateUrl: './upload.html'
})
export class FileUploadComponent  {
    uploadedFiles: any[] = [];
    private patientId: string;

    constructor(
        private route: ActivatedRoute,
        private fileService: FileService) {
        this.patientId = this.route.snapshot.params['id'];
    }

    myUploader($event) {
        $event.files.forEach(f=>this.postFile(f));
    }

    postFile(fileToUpload: File): void {
        this.fileService.postFile(this.patientId, fileToUpload).pipe(take(1)).subscribe();
    }
}