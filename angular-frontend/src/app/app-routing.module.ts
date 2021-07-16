import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { AppComponent } from './app.component';
import { LoginComponent } from './login/login.component';
import { ProtectedComponent } from './protected/protected.component';
import {CallbackComponent} from "./callback/callback.component";
import {AutoLoginPartialRoutesGuard} from "angular-auth-oidc-client";

const routes: Routes = [
  {path: 'login', component: LoginComponent},
  {path: 'protected', component: ProtectedComponent,  canActivate: [AutoLoginPartialRoutesGuard]},
  {path: 'callback', component: CallbackComponent},
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
