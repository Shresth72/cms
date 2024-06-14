export default function Home() {
    return (
        <main className="bg-gray-900 h-screen">

        {/* <!-- Navbar --> */}
        <nav className="bg-dark text-white shadow">
          <div className="container mx-auto px-6 py-3 flex justify-between items-center">
            <div className="flex items-center">
              {/* <img src="image.png" alt="Landscape Finance Lab" className="h-10 w-10 mr-3"/> */}
              <span className="text-xl font-semibold">LANDSCAPE FINANCE LAB</span>
            </div>
            <div className="hidden md:flex space-x-4">
              <a href="#" className="hover:text-gray-400">Landscapes</a>
              <a href="#" className="hover:text-gray-400">Finance</a>
              <a href="#" className="hover:text-gray-400">Learning</a>
              <a href="#" className="hover:text-gray-400">About</a>
              <a href="#" className="hover:text-gray-400">Updates</a>
              <a href="#" className="hover:text-gray-400">Contact</a>
            </div>
            <div className="md:hidden">
              <button id="menu-btn" className="text-white focus:outline-none">
                <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16m-7 6h7"></path>
                </svg>
              </button>
            </div>
          </div>
          <div id="menu" className="md:hidden">
            <a href="#" className="block px-4 py-2 text-sm text-white hover:bg-gray-700">Landscapes</a>
            <a href="#" className="block px-4 py-2 text-sm text-white hover:bg-gray-700">Finance</a>
            <a href="#" className="block px-4 py-2 text-sm text-white hover:bg-gray-700">Learning</a>
            <a href="#" className="block px-4 py-2 text-sm text-white hover:bg-gray-700">About</a>
            <a href="#" className="block px-4 py-2 text-sm text-white hover:bg-gray-700">Updates</a>
            <a href="#" className="block px-4 py-2 text-sm text-white hover:bg-gray-700">Contact</a>
          </div>
        </nav>
      
        {/* <!-- Impacts Section --> */}
        <section className="bg-dark text-white py-16 ">
          <div className="container mx-auto text-center">
            <h2 className="text-2xl font-semibold mb-15">Our Impacts</h2>
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-2 gap-8">
              <div>
                <p className="text-5xl font-bold mb-2">85M+</p>
                <p className="text-lg">Hectares under restoration with Lab support</p>
              </div>
              <div>
                <p className="text-5xl font-bold mb-2">€15M+</p>
                <p className="text-lg">Directly-secured investments into landscapes</p>
              </div>
              <div>
                <p className="text-5xl font-bold mb-2">€7M+</p>
                <p className="text-lg">Pipeline of landscape opportunities under development</p>
              </div>
              <div>
                <p className="text-5xl font-bold mb-2">2000+</p>
                <p className="text-lg">Landscape practitioners & investors in our learning networks</p>
              </div>
            </div>
          </div>
        </section>
      
        {/* <script>
          document.getElementById('menu-btn').addEventListener('click', function () {
            var menu = document.getElementById('menu');
            menu.classNameList.toggle('hidden');
          });
        </script> */}
      
      </main>
    );
}