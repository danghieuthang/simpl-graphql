import { Component, OnInit, WritableSignal } from '@angular/core';
import { Apollo } from 'apollo-angular';
import { GET_USERS } from '../graphql/graphql.queries';
import { takeUntil } from 'rxjs';

export interface IUser {
  id: number;
  name: string;
  email: string;
  slug: string;
}

@Component({
  selector: 'app-user-list',
  templateUrl: './user-list.component.html',
  styleUrls: ['./user-list.component.scss'],
}) export class UserListComponent implements OnInit {
  users: IUser[] = [];
  error?: string;

  constructor(private _appolo: Apollo) {
  }
  randomData() {

    for (let index = 0; index < 5; index++) {
      this.users.push({
        id: index,
        name: "Name " + index,
        email: "email" + index + "@gmail.com",
        slug: "user-" + index
      })
    }
  }
  getUsers() {
    this._appolo.watchQuery({
      query: GET_USERS
    })
      .valueChanges.subscribe({
        next: ({ data, error: errors }: any) => {
          this.users = data.users.data.map((u: any) => {
            const user = {
              id: u.id,
              name: u.name,
              email: u.email,
              slug: u.id
            }
            return user
          });
          this.error = errors?.message
        }, complete: () => console.log('complete')
      })
  }
  ngOnInit(): void {
    // this.randomData()
    this.getUsers()
  }
}
