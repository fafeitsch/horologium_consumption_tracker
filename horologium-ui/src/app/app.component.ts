import {Component, OnInit} from '@angular/core';
import gql from 'graphql-tag';
import {Apollo} from 'apollo-angular';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.sass']
})
export class AppComponent  implements OnInit {
  title = 'horologium-ui';

  constructor(private apollo: Apollo) { }

  ngOnInit(): void {
    this.apollo.query({
      query: gql`query {
        allSeries{name}
      }`
    }).subscribe(data => {
      console.log(data);
    });
  }


}
