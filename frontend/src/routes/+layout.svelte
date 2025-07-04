<script>
	import { page } from '$app/stores';
	import { browser } from '$app/environment';
	import { onMount } from 'svelte';
	import '../app.css';

	// Theme management
	let theme = 'light';
	let sidebarOpen = false;

	onMount(() => {
		// Load theme from localStorage
		if (browser) {
			const savedTheme = localStorage.getItem('theme') || 'light';
			theme = savedTheme;
			document.documentElement.setAttribute('data-theme', theme);
		}
	});

	function toggleTheme() {
		theme = theme === 'light' ? 'dark' : 'light';
		
		if (browser) {
			document.documentElement.setAttribute('data-theme', theme);
			localStorage.setItem('theme', theme);
		}
	}

	function toggleSidebar() {
		sidebarOpen = !sidebarOpen;
	}

	// Navigation items
	const navItems = [
		{
			section: 'Content',
			items: [
				{
					label: 'Feeds',
					href: '/feeds',
					icon: '📡',
					active: false
				},
				{
					label: 'Articles',
					href: '/articles',
					icon: '📄',
					active: false
				}
			]
		},
		{
			section: 'Organization',
			items: [
				{
					label: 'Lists',
					href: '/lists',
					icon: '📋',
					disabled: true,
					badge: 'Soon'
				}
			]
		}
	];

	// Update active states based on current page
	$: {
		navItems.forEach(section => {
			section.items.forEach(item => {
				item.active = $page.url.pathname.startsWith(item.href) && item.href !== '/';
			});
		});
	}
</script>

<svelte:head>
	<title>RSS Reader</title>
	<meta name="description" content="A modern RSS feed reader" />
</svelte:head>

<div class="app-container">
	<!-- Header -->
	<header class="header">
		<a href="/" class="header-title">
			RSS Reader
		</a>
		
		<div class="header-actions">
			<!-- Mobile menu toggle -->
			<button 
				class="btn-ghost md:hidden"
				on:click={toggleSidebar}
				aria-label="Toggle sidebar"
			>
				☰
			</button>
			
			<!-- Theme toggle -->
			<button 
				class="theme-toggle"
				on:click={toggleTheme}
				aria-label="Toggle theme"
				title="Toggle {theme === 'light' ? 'dark' : 'light'} mode"
			>
				{theme === 'light' ? '🌙' : '☀️'}
			</button>
		</div>
	</header>

	<!-- Sidebar -->
	<aside class="sidebar {sidebarOpen ? 'open' : ''}">
		<div class="sidebar-header">
			<h2 class="sidebar-title">Navigation</h2>
		</div>
		
		<nav class="sidebar-nav">
			{#each navItems as section}
				<div class="nav-section">
					<h3 class="nav-section-title">{section.section}</h3>
					<ul class="nav-list">
						{#each section.items as item}
							<li class="nav-item">
								<a
									href={item.disabled ? '#' : item.href}
									class="nav-link {item.active ? 'active' : ''} {item.disabled ? 'disabled' : ''}"
									on:click={() => sidebarOpen = false}
								>
									<span class="nav-icon">{item.icon}</span>
									<span class="nav-text">{item.label}</span>
									{#if item.badge}
										<span class="nav-badge">{item.badge}</span>
									{/if}
								</a>
							</li>
						{/each}
					</ul>
				</div>
			{/each}
		</nav>
	</aside>

	<!-- Main Content -->
	<main class="main-content">
		<slot />
	</main>
</div>

<!-- Mobile sidebar backdrop -->
{#if sidebarOpen}
	<button 
		class="fixed inset-0 bg-black bg-opacity-50 z-10 md:hidden backdrop-button"
		on:click={toggleSidebar}
		aria-label="Close sidebar"
	></button>
{/if}

<style>
	/* Additional responsive utilities */
	@media (min-width: 768px) {
		.md\:hidden {
			display: none;
		}
	}
	
	.fixed {
		position: fixed;
	}
	
	.inset-0 {
		top: 0;
		right: 0;
		bottom: 0;
		left: 0;
	}
	
	.bg-black {
		background-color: rgb(0 0 0);
	}
	
	.bg-opacity-50 {
		background-color: rgb(0 0 0 / 0.5);
	}
	
	.z-10 {
		z-index: 10;
	}
</style>