import { Component, OnInit } from '@angular/core';
import { FileService, IFileModel } from '../../services/file.service';
import { ActivatedRoute } from '@angular/router';
import { NotificationsManager } from '../../services/notifications-manager';
import { PatientService } from '../../services/patient.service';
import { Observable } from 'rxjs';
import { take } from 'rxjs/operators';

@Component({
    selector: 'download',
    templateUrl: './download.html'
})
export class DownloadComponent implements OnInit  {

    private patientId: string;
    tags: Observable<string[]>;

    selectedTags: string[] = [];

    files: IFileModel[];

    constructor(
        private route: ActivatedRoute,
        private fileService: FileService,
        private patientService: PatientService,
        private notificationsManager: NotificationsManager) {

        this.patientId = this.route.snapshot.params['id'];
        this.tags = this.patientService.getTagsPerPatient(this.patientId);
    }

    ngOnInit(): void {
    }

    selectedTagsChanged() {
        this.fileService.getFilesByTags(this.patientId, this.selectedTags).pipe(take(1)).subscribe(n => this.files = n);
    }

    download($event: any, file: IFileModel ) {
        this.fileService.downloadFile(this.patientId, file);
    }
}

