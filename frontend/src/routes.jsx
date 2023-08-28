// pages
import Login from './pages/Login/Login';
import PanelPrivateRoute from './components/Privates/PanelPrivateRoute';
import Panel from './pages/Panel/Panel';

// components
import PeerDetails from './components/PeerDetails/PeerDetails';
import NewPeer from './components/NewPeer/NewPeer';

// routes
const routes = [
	{
		path: '/',
		element: (
			<PanelPrivateRoute>
				<Panel />
			</PanelPrivateRoute>
		),
		children: [
			{ path: 'new-peer', element: <NewPeer /> },
			{ path: 'peer', element: <PeerDetails /> }
		]
	},
	{ path: '/login', element: <Login /> }
];

// exports
export default routes;
