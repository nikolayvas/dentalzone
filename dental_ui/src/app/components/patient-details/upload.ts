import { Component, ViewChild } from '@angular/core';
import { FileService } from '../../services/file.service';
import { ActivatedRoute } from '@angular/router';
import { take } from 'rxjs/operators';
import { NotificationsManager } from '../../services/notifications-manager';
import { FormGroup, FormControl } from '@angular/forms';

@Component({
    selector: 'upload',
    templateUrl: './upload.html'
})
export class FileUploadComponent  {
    private patientId: string;

    @ViewChild('uploader') uploader;

    form: FormGroup;
    tags: FormControl;

    constructor(
        private route: ActivatedRoute,
        private fileService: FileService,
        private notificationsManager: NotificationsManager) {
        this.patientId = this.route.snapshot.params['id'];

        this.form = new FormGroup({
            tags: this.tags = new FormControl(''),
        });
    }

    myUploader($event) {
        $event.files.forEach(f=>this.postFile(f));

        this.notificationsManager.ServerSuccess("Selected files were successfully uploaded.");
        this.uploader.clear();
    }

    postFile(fileToUpload: File): void {
        var tags = (this.tags.value as string).split(/\s+/);
        this.fileService.postFile(this.patientId, fileToUpload, tags).pipe(take(1)).subscribe();
    }
}