import React, { Component, PropTypes } from 'react';
import { connect } from 'react-redux';
import { performLogin } from '../actions';
import { Link } from 'react-router';

export default class Dashboard extends Component {
  render() {
    return (
      <div className="container">
        <h2>Dashboard</h2>
        <div className="row metrics-row">
          <div className="col-sm-4 metrics-slot metrics-slot-green">
            <h1 className="metrics-slot-jumbonum">1</h1>
            <strong>healthy</strong> host
          </div>
          <div className="col-sm-4 metrics-slot metrics-slot-yellow">
            <h1 className="metrics-slot-jumbonum">0</h1>
            <strong>warning</strong> hosts
          </div>
          <div className="col-sm-4 metrics-slot metrics-slot-red">
            <h1 className="metrics-slot-jumbonum">0</h1>
            <strong>alerting</strong> hosts
          </div>
        </div>


        <div className="row metrics-row">
          <h3>Recent Events</h3>
          <ul className="logs-list">
            <li>[2016-02-21T04:13:52Z] <Link to="/hosts/2">(lin.lukegb.com)</Link> transaction /717_ecdbebbc from uid 1000 finished with success after 0ms</li>
            <li>[2016-02-21T04:13:52Z] <Link to="/hosts/2">(lin.lukegb.com)</Link> transaction /716_bcbbadbe from uid 1000 finished with success after 163ms</li>
            <li>[2016-02-21T04:13:52Z] <Link to="/hosts/2">(lin.lukegb.com)</Link> transaction /715_eaeddbca from uid 1000 finished with success after 610ms</li>
            <li>[2016-02-21T04:13:52Z] <Link to="/hosts/2">(lin.lukegb.com)</Link> 9896): PackageKit-alpm-WARNING **: failed retrieving file 'multilib.db' from archlinux.mirrors.house.lukegb.com : Could not resolve host: archlinux.mirrors.house.lukegb.com</li>
            <li>[2016-02-21T04:13:52Z] <Link to="/hosts/2">(lin.lukegb.com)</Link> 9896): PackageKit-alpm-WARNING **: failed retrieving file 'multilib-testing.db' from archlinux.mirrors.house.lukegb.com : Could not resolve host: archlinux.mirrors.house.lukegb.com</li>
            <li>[2016-02-21T04:13:52Z] <Link to="/hosts/2">(lin.lukegb.com)</Link> 9896): PackageKit-alpm-WARNING **: failed retrieving file 'community.db' from archlinux.mirrors.house.lukegb.com : Could not resolve host: archlinux.mirrors.house.lukegb.com</li>
            <li>[2016-02-21T04:13:52Z] <Link to="/hosts/2">(lin.lukegb.com)</Link> 9896): PackageKit-alpm-WARNING **: failed retrieving file 'community-testing.db' from archlinux.mirrors.house.lukegb.com : Could not resolve host: archlinux.mirrors.house.lukegb.com</li>
            <li>[2016-02-21T04:13:52Z] <Link to="/hosts/2">(lin.lukegb.com)</Link> 9896): PackageKit-alpm-WARNING **: failed retrieving file 'extra.db' from archlinux.mirrors.house.lukegb.com : Could not resolve host: archlinux.mirrors.house.lukegb.com</li>

          </ul>

        </div>

      </div>
    );
  }
}
