import './App.css';
import icon from '../../build/windows/icon.ico'
import { Loader } from './shared/Loader/Loader';
import { useEffect } from 'react';

function App() {
    useEffect(() => {
        const onWheel = (e: WheelEvent) => {
            if (e.ctrlKey) {
                e.preventDefault();
            }
        };

        const onKeyDown = (e: KeyboardEvent) => {
            if (e.ctrlKey && ['+', '-', '=', '0'].includes(e.key)) {
                e.preventDefault();
            }
        };

        window.addEventListener('wheel', onWheel, { passive: false });
        window.addEventListener('keydown', onKeyDown);

        return () => {
            window.removeEventListener('wheel', onWheel);
            window.removeEventListener('keydown', onKeyDown);
        };
    }, []);


    return (
        <div id="App">
            <img className='icon' src={icon}></img>
            <Loader />
        </div>
    )
}

export default App
