import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { Apollo } from 'apollo-angular';
import { GET_USER_DETAIL } from '../graphql/graphql.queries';
import { IUser } from '../user-list/user-list.component';

interface RoleDetail{
  name: string
}
interface UserDetail{
  id: number;
  name: string;
  email: string;
  createdAt: Date;
  lastModifiedAt: Date;
  role?: RoleDetail;
}

@Component({
  selector: 'app-user-detail',
  templateUrl: './user-detail.component.html',
  styleUrls: ['./user-detail.component.scss']
})
export class UserDetailComponent implements OnInit{
  user?: UserDetail
  error?:string
  constructor(private _route: ActivatedRoute, private _appolo: Apollo){

  }
  ngOnInit(): void {
    let id = this._route.snapshot.paramMap.get('slug')
    this._appolo.watchQuery({
      query: GET_USER_DETAIL,
      variables:{
        id:id
      }
    }).valueChanges.subscribe(({data,errors}: any)=>{
        this.user = data.user as UserDetail
        this.error = errors[0]?.message
    })
  }
  
}
