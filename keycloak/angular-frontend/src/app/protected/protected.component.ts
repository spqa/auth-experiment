import { Component, OnInit } from '@angular/core';

@Component({
  selector: 'app-protected',
  templateUrl: './protected.component.html',
  styleUrls: ['./protected.component.scss']
})
export class ProtectedComponent implements OnInit {

  profile: any;

  constructor() {
   }

  ngOnInit(): void {

  }

}
