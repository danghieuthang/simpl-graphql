import { gql } from "apollo-angular";

const GET_USERS = gql`
query{
  users(page: {page: 1, pageSize: 10}) {
    data {
      id
      name
      email
      createdAt
      lastModifiedAt
      role {
        id
        name
      }
    }
    pageInfo {
      total
      page
      pageSize
    }
  }
}

`

const GET_USER_DETAIL = gql`
query user($id: Int!) {
  user(id: $id) {
    id
    name
    lastModifiedAt
    email
    createdBy
    createdAt
    updatedBy
    lastModifiedAt
    role {
      id
      name
    }
  }
}
`

export { GET_USERS, GET_USER_DETAIL }