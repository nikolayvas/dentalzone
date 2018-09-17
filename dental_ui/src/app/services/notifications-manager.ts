
import {throwError as observableThrowError,  Observable } from 'rxjs';
import { Injectable } from '@angular/core';
import { MessageService } from 'primeng/api';
import { HttpErrorResponse } from '@angular/common/http';

@Injectable()
export class NotificationsManager{

  constructor(private messageService: MessageService){

  }

  public ServerError(error: any, customMessage?: string) {

    if (error instanceof HttpErrorResponse) {
      let errMessage = '';
      try {
        errMessage = error.message;
        const errorMsg = "Server Request (" + customMessage+ "), Resposne: " + errMessage;
        this.messageService.add({severity:'error', summary: 'Server ERROR:', detail: errorMsg });
      } catch( err ) {
        errMessage = error.statusText;
        const errorMsg = "Server Request (" + customMessage + ") -  Status Text: " + errMessage;
        this.messageService.add({severity:'error', summary: 'Server ERROR:', detail: errorMsg});
      }

      return observableThrowError(errMessage);
    }
  }

  public ServerSuccess(message: string, title?: string){
    this.messageService.add({severity:'success', summary: title, detail: message});
  }

  public WarningNotification(message: string){
    this.messageService.add({severity:'warn', summary: 'Warning', detail: message});
  }
}
